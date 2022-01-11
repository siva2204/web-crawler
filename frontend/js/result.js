const getData = (data) => {
	const highest = 1.01 * data[0].rank;
	return data.map((d) => {
		return `<div class="item">
		<div class="link">${d.url}</div>
		<div class="title">
			<a href="${d.url}">
			${d.title ? d.title : d.url}
			</a>
			<span class="match">(${
				Math.round((d.rank / highest) * 100 * 100) / 100
			}% match)</span>
		</div>
		<div class="descp">${d.description}</div>
	</div>`;
	});
};

const start = async () => {
	const urlSearchParams = new URLSearchParams(window.location.search);
	const params = Object.fromEntries(urlSearchParams.entries());

	document.getElementById("search-input").value = params.search;

	let timeStart = Date.now();
	let resp = await axios.get("/search/" + params.search);
	let timeEnd = Date.now();
	let { data } = resp;
	data = data.data;
	console.log(data);

	data.sort((a, b) => a.rank < b.rank);

	document.getElementById("main-data").innerHTML = getData(data);
	document.getElementById("time-taken").innerHTML = `Found ${
		data.length
	} results in ${(timeEnd - timeStart) / 1000} s`;
};
window.addEventListener("load", (event) => {
	console.log("page is fully loaded");
	start();
	autocomplete(document.getElementById("search-input"));
	loadNewPage();
});

function autocomplete(inp) {
	var currentFocus;
	inp.addEventListener("input", async function (e) {
		var a,
			b,
			i,
			val = this.value;
		let { data } = await axios.post("/words/" + val);
		data = data.data;
		arr = data.slice(0, 30);
		closeAllLists();
		if (!val) {
			return false;
		}
		currentFocus = -1;
		a = document.createElement("DIV");
		a.setAttribute("id", this.id + "autocomplete-list");
		a.setAttribute("class", "autocomplete-items");
		console.log(a);
		this.parentNode.appendChild(a);
		for (i = 0; i < arr.length; i++) {
			if (arr[i].substr(0, val.length).toUpperCase() == val.toUpperCase()) {
				b = document.createElement("DIV");
				b.innerHTML = "<strong>" + arr[i].substr(0, val.length) + "</strong>";
				b.innerHTML += arr[i].substr(val.length);
				b.innerHTML += "<input type='hidden' value='" + arr[i] + "'>";
				b.addEventListener("click", function (e) {
					inp.value = this.getElementsByTagName("input")[0].value;
					closeAllLists();
				});
				a.appendChild(b);
			}
		}
	});
	inp.addEventListener("keydown", function (e) {
		var x = document.getElementById(this.id + "autocomplete-list");
		if (x) x = x.getElementsByTagName("div");
		if (e.keyCode == 40) {
			currentFocus++;
			addActive(x);
		} else if (e.keyCode == 38) {
			currentFocus--;
			addActive(x);
		} else if (e.keyCode == 13) {
			e.preventDefault();
			if (currentFocus > -1) {
				if (x) x[currentFocus].click();
			}
		}
	});

	function addActive(x) {
		if (!x) return false;
		removeActive(x);
		if (currentFocus >= x.length) currentFocus = 0;
		if (currentFocus < 0) currentFocus = x.length - 1;
		x[currentFocus].classList.add("autocomplete-active");
	}

	function removeActive(x) {
		for (var i = 0; i < x.length; i++) {
			x[i].classList.remove("autocomplete-active");
		}
	}

	function closeAllLists(elmnt) {
		var x = document.getElementsByClassName("autocomplete-items");
		for (var i = 0; i < x.length; i++) {
			if (elmnt != x[i] && elmnt != inp) {
				x[i].parentNode.removeChild(x[i]);
			}
		}
	}
	document.addEventListener("click", function (e) {
		closeAllLists(e.target);
	});
}

function loadNewPage() {
	document.getElementById("search-form").addEventListener("submit", (e) => {
		e.preventDefault();
		const value = document.getElementById("search-input").value;
		const a = document.createElement("a");
		a.href = "/result?search=" + value;
		a.click();
	});
}
