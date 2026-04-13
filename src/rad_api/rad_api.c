#define _CRT_SECURE_NO_WARNINGS
#include <windows.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <tlhelp32.h>

#define IPC_SHARED_MEMORY_BUFFER_SIZE (4u * 1024u * 1024u)
#define MAX_LINE 65536

typedef struct IPCInfo {
  uint64_t msg_size;
} IPCInfo;

typedef struct {
  DWORD pid;
  HANDLE s2m_map;
  HANDLE s2m_signal;
  HANDLE s2m_lock;
  HANDLE m2s_map;
  HANDLE m2s_signal;
  HANDLE m2s_lock;
  unsigned char *s2m_base;
  unsigned char *m2s_base;
} RadApiState;

static RadApiState RAD = {0};

static void build_ipc_name(char *dst, size_t dst_cap, const char *kind,
                           DWORD pid) {
  _snprintf(dst, dst_cap, "_raddbg_ipc_%s_%lu_", kind, (unsigned long)pid);
  dst[dst_cap - 1] = 0;
}

static DWORD find_raddbg_pid(void) {
  DWORD result = 0;
  DWORD self_pid = GetCurrentProcessId();
  HANDLE snap = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
  if (snap == INVALID_HANDLE_VALUE) {
    return 0;
  }

  PROCESSENTRY32 pe = {0};
  pe.dwSize = sizeof(pe);
  if (Process32First(snap, &pe)) {
    do {
      if (pe.th32ProcessID != self_pid &&
          (_stricmp(pe.szExeFile, "raddbg.exe") == 0 ||
           _stricmp(pe.szExeFile, "raddbg") == 0)) {
        result = pe.th32ProcessID;
        break;
      }
    } while (Process32Next(snap, &pe));
  }

  CloseHandle(snap);
  return result;
}

static void RadRelease() {
  if (RAD.s2m_base)
    UnmapViewOfFile(RAD.s2m_base);
  if (RAD.m2s_base)
    UnmapViewOfFile(RAD.m2s_base);
  if (RAD.s2m_map)
    CloseHandle(RAD.s2m_map);
  if (RAD.m2s_map)
    CloseHandle(RAD.m2s_map);
  if (RAD.s2m_signal)
    CloseHandle(RAD.s2m_signal);
  if (RAD.s2m_lock)
    CloseHandle(RAD.s2m_lock);
  if (RAD.m2s_signal)
    CloseHandle(RAD.m2s_signal);
  if (RAD.m2s_lock)
    CloseHandle(RAD.m2s_lock);
  memset(&RAD, 0, sizeof(RAD));
}

static int ipc_open(RadApiState *RAD) {
  char name[256];

  build_ipc_name(name, sizeof(name), "sender2main_shared_memory", RAD->pid);
  RAD->s2m_map = OpenFileMappingA(FILE_MAP_READ | FILE_MAP_WRITE, FALSE, name);
  build_ipc_name(name, sizeof(name), "sender2main_signal_semaphore", RAD->pid);
  RAD->s2m_signal =
      OpenSemaphoreA(SYNCHRONIZE | SEMAPHORE_MODIFY_STATE, FALSE, name);
  build_ipc_name(name, sizeof(name), "sender2main_lock_semaphore", RAD->pid);
  RAD->s2m_lock =
      OpenSemaphoreA(SYNCHRONIZE | SEMAPHORE_MODIFY_STATE, FALSE, name);

  build_ipc_name(name, sizeof(name), "main2sender_shared_memory", RAD->pid);
  RAD->m2s_map = OpenFileMappingA(FILE_MAP_READ | FILE_MAP_WRITE, FALSE, name);
  build_ipc_name(name, sizeof(name), "main2sender_signal_semaphore", RAD->pid);
  RAD->m2s_signal =
      OpenSemaphoreA(SYNCHRONIZE | SEMAPHORE_MODIFY_STATE, FALSE, name);
  build_ipc_name(name, sizeof(name), "main2sender_lock_semaphore", RAD->pid);
  RAD->m2s_lock =
      OpenSemaphoreA(SYNCHRONIZE | SEMAPHORE_MODIFY_STATE, FALSE, name);

  if (!RAD->s2m_map || !RAD->s2m_signal || !RAD->s2m_lock || !RAD->m2s_map ||
      !RAD->m2s_signal || !RAD->m2s_lock) {
    RadRelease();
    return 0;
  }

  RAD->s2m_base = (unsigned char *)MapViewOfFile(
      RAD->s2m_map, FILE_MAP_READ | FILE_MAP_WRITE, 0, 0,
      IPC_SHARED_MEMORY_BUFFER_SIZE);
  RAD->m2s_base = (unsigned char *)MapViewOfFile(
      RAD->m2s_map, FILE_MAP_READ | FILE_MAP_WRITE, 0, 0,
      IPC_SHARED_MEMORY_BUFFER_SIZE);

  if (!RAD->s2m_base || !RAD->m2s_base) {
    RadRelease();
    return 0;
  }

  return 1;
}

static int ipc_send_and_wait(RadApiState *RAD, const char *text,
                             size_t text_size) {
  const size_t max_payload = IPC_SHARED_MEMORY_BUFFER_SIZE - sizeof(IPCInfo);
  if (text_size > max_payload) {
    text_size = max_payload;
  }

  DWORD wait = WaitForSingleObject(RAD->s2m_lock, INFINITE);
  if (wait != WAIT_OBJECT_0) {
    fprintf(stderr, "failed waiting sender2main lock\n");
    return 4;
  }

  {
    IPCInfo *info = (IPCInfo *)RAD->s2m_base;
    unsigned char *payload = (unsigned char *)(info + 1);
    info->msg_size = (uint64_t)text_size;
    if (text_size != 0) {
      memcpy(payload, text, text_size);
    }
  }

  ReleaseSemaphore(RAD->s2m_signal, 1, NULL);
  ReleaseSemaphore(RAD->s2m_lock, 1, NULL);

  wait = WaitForSingleObject(RAD->m2s_signal, 10000);
  if (wait != WAIT_OBJECT_0) {
    fprintf(stderr, "timeout/failure waiting reply from raddbg\n");
    return 5;
  }

  wait = WaitForSingleObject(RAD->m2s_lock, INFINITE);
  if (wait != WAIT_OBJECT_0) {
    fprintf(stderr, "failed waiting main2sender lock\n");
    return 6;
  }

  {
    IPCInfo *info = (IPCInfo *)RAD->m2s_base;
    unsigned char *payload = (unsigned char *)(info + 1);
    size_t size = (size_t)info->msg_size;
    size_t i = 0;
    if (size > max_payload) {
      size = max_payload;
    }
    while (i < size) {
      size_t start = i;
      while (i < size && payload[i] != 0) {
        i += 1;
      }
      if (i > start) {
        fwrite(payload + start, 1, i - start, stdout);
      }
      i += 1;
    }
    fflush(stdout);
  }

  ReleaseSemaphore(RAD->m2s_lock, 1, NULL);
  return 0;
}

static int RadSendCommand(const char *line,
                            size_t line_size) {
  const char *prefix = "--ipc ";
  size_t prefix_size = strlen(prefix);
  size_t total = prefix_size + line_size;
  char *msg = (char *)malloc(total + 1);
  int ok = 0;
  if (msg == NULL) {
    fprintf(stderr, "out of memory\n");
    return 3;
  }
  memcpy(msg, prefix, prefix_size);
  memcpy(msg + prefix_size, line, line_size);
  msg[total] = 0;
  ok = ipc_send_and_wait(&RAD, msg, total);
  free(msg);

  return ok;
}

static int RadInit() {
  RAD.pid = find_raddbg_pid();
  if (RAD.pid == 0) {
    return 1;
  }

  if (!ipc_open(&RAD)) {
    fprintf(stderr, "Failed to open raddbg IPC resources for pid %lu.\n",
            (unsigned long)RAD.pid);
    return 2;
  }

  printf("Connected to raddbg pid %lu\n", (unsigned long)RAD.pid);

  return 0;
}
