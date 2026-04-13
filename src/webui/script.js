window.addRow = function (id, functionName) {
	const tbody = document.querySelector("#funcTable tbody");
	if (!tbody) return;

	const rowCount = tbody.rows.length;

	const tr = document.createElement("tr");

	const tdSrNo = document.createElement("td");
	tdSrNo.textContent = rowCount + 1;
	const tdID = document.createElement("td");
	tdID.textContent = id;

	const tdFunc = document.createElement("td");
	tdFunc.textContent = functionName;

	const tdTime = document.createElement("td");
	const now = new Date();
	const hours = String(now.getHours()).padStart(2, "0");
	const minutes = String(now.getMinutes()).padStart(2, "0");
	const seconds = String(now.getSeconds()).padStart(2, "0");
	tdTime.textContent = `${hours}:${minutes}:${seconds}`;

	tr.onclick = function () {
		loadDescriptionFromGo(id);
	};

	tr.appendChild(tdSrNo);
	tr.appendChild(tdTime);
	tr.appendChild(tdID);
	tr.appendChild(tdFunc);

	tbody.appendChild(tr);
};

(function () {
	const divider = document.getElementById("divider");
	const tablePanel = document.getElementById("tablePanel");
	const container = document.getElementById("splitContainer");

	let isDragging = false;

	divider.addEventListener("mousedown", function () {
		isDragging = true;
		document.body.style.cursor = "col-resize";
		document.body.style.userSelect = "none";
	});

	document.addEventListener("mousemove", function (e) {
		if (!isDragging) return;

		const rect = container.getBoundingClientRect();
		const offsetX = e.clientX - rect.left;

		let newWidth = offsetX;

		tablePanel.style.width = newWidth + "px";
	});

	document.addEventListener("mouseup", function () {
		isDragging = false;
		document.body.style.cursor = "";
		document.body.style.userSelect = "";
	});
})();

window.loadDescriptionFromGo = async function (id) {
	try {
		const html = await getDescFromGo(id);

		const descContent = document.getElementById("descContent");
		if (!descContent) return;

		descContent.innerHTML = "<pre>" + html + "</pre>";
	} catch (err) {
		console.error("Failed to load description from Go:", err);
	}
};

document.addEventListener("click", function (e) {
	const link = e.target.closest("a");
	if (!link) return;

	const href = link.getAttribute("href");
	if (!href) return;

	if (!href.startsWith("http://") && !href.startsWith("https://")) {
		return;
	}

	e.preventDefault();

	openExternalLink(href)
		.then(res => console.log("Opened link:", href, res))
		.catch(err => console.error("Failed to open link:", err));
});

function searchFunction() {
	const searchInput = document.getElementById("searchBar");
	const filter = searchInput.value.trim().toLowerCase();

	const table = document.getElementById("funcTable");
	const rows = table.getElementsByTagName("tr");

	for (let i = 1; i < rows.length; i++) {
		rows[i].style.backgroundColor = "";
		rows[i].style.color = "#ffffff"
	}

	if (filter === "") return;

	for (let i = 1; i < rows.length; i++) {
		const cells = rows[i].getElementsByTagName("td");

		if (cells.length >= 3) {
			const functionName = cells[3].textContent.toLowerCase();

			if (functionName.includes(filter)) {
				rows[i].style.backgroundColor = "#baa87cff";
				rows[i].style.color = "#000000";
			}
		}
	}
}

function handleClick(action) {
	console.log("Clicked:", action);

	debugCommand(action)
		.then(() => {
			console.log("Command sent:", action);
		})
		.catch(err => {
			console.error("Error:", err);
		});
}

document.addEventListener("click", function (event) {
	const panel = document.getElementById("injectPanel");
	const injectBtn = document.querySelector(".inject-btn");

	if (panel.style.display !== "block") return;

	const isInsidePanel = panel.contains(event.target);

	const isInjectButton = injectBtn.contains(event.target);

	if (!isInsidePanel && !isInjectButton) {
		panel.style.display = "none";
	}
});

function browseFile(type) {
	// Call Go
	console.log("searching for:", type);

	openFileDialog(type).then(path => {
		if (!path) return;

		if (type === "dll") {
			document.getElementById("dllPath").value = path;
		} else if (type === "exe") {
			document.getElementById("exePath").value = path;
		}
	});
}

function toggleInjectPanel() {
	const panel = document.getElementById("injectPanel");

	if (panel.style.display === "block") {
		panel.style.display = "none";
	} else {
		panel.style.display = "block";
	}
}

function handleInject(remove = false) {
	const dllPath = document.getElementById("dllPath").value.trim();
	const exePath = document.getElementById("exePath").value.trim();

	if (!dllPath) {
		alert("Please select DLL path");
		return;
	}

	if (!exePath) {
		alert("Please select executable path");
		return;
	}

	console.log("DLL:", dllPath);
	console.log("EXE:", exePath);

	injectHookDll(dllPath, exePath, remove)
		.then(() => {
			console.log("Injection success");
			alert("DLL Injected Successfully");
		})
		.catch(err => {
			console.error("Injection failed:", err);
			alert("Injection failed: " + err);
		});
}

window.addEventListener("load", function () {
	if (typeof onPageReload === "function") {
		onPageReload()
			.then(res => console.log("Go returned:", res))
			.catch(err => console.error(err));
	}
});