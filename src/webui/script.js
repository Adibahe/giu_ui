window.addRow = function (id, functionName) {
	const tbody = document.querySelector("table tbody");
	if (!tbody) return;

	const rowCount = tbody.rows.length;

	const tr = document.createElement("tr");

	const tdSrNo = document.createElement("td");
	tdSrNo.textContent = rowCount + 1;

	const tdID = document.createElement("td");
	tdID.textContent = id;

	const tdFunc = document.createElement("td");
	tdFunc.textContent = functionName;

	tr.onclick = function () {
		loadDescriptionFromGo(id);
	};

	tr.appendChild(tdSrNo);
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
