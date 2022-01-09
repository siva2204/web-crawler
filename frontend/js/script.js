function autocomplete(inp) {
    var currentFocus;
    inp.addEventListener("input", async function(e) {
        var a,
            b,
            i,
            val = this.value;
        let {
            data
        } = await axios.post("/words/" + val);
        data = data.data
        arr = data.slice(0, 30);
        closeAllLists();
        if (!val) {
            return false;
        }
        currentFocus = -1;
        a = document.createElement("DIV");
        a.setAttribute("id", this.id + "autocomplete-list");
        a.setAttribute("class", "autocomplete-items");
        this.parentNode.appendChild(a);
        for (i = 0; i < arr.length; i++) {
            if (
                arr[i].substr(0, val.length).toUpperCase() == val.toUpperCase()
            ) {
                b = document.createElement("DIV");
                b.innerHTML =
                    "<strong>" + arr[i].substr(0, val.length) + "</strong>";
                b.innerHTML += arr[i].substr(val.length);
                b.innerHTML += "<input type='hidden' value='" + arr[i] + "'>";
                b.addEventListener("click", function(e) {
                    inp.value = this.getElementsByTagName("input")[0].value;
                    closeAllLists();
                });
                a.appendChild(b);
            }
        }
    });
    inp.addEventListener("keydown", function(e) {
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
    document.addEventListener("click", function(e) {
        closeAllLists(e.target);
    });
}
autocomplete(document.getElementById("input"));

document.querySelector("#search-button").onclick = async(e) => {
    e.preventDefault();
    word = document.querySelector("#input").value;
    document.querySelector("section").innerText = "Searching...";
    let timeStart = Date.now();
    let {
        data
    } = await axios.get("/search/" + word);

    data = JSON.parse(data.data)['Payload'];
    data = [...new Set(data)];
    let timeEnd = Date.now();
    let ele = document.querySelector("#links");
    ele.innerHTML = "";
    for (var i = 0; i < data.length; i++) {
        var aside = document.createElement("aside");
        var a = document.createElement("a");
        var link = document.createTextNode(`${data[i]}`);
        a.appendChild(link);
        a.target = "_blank";
        a.href = `${data[i]}`;
        aside.appendChild(a);
        ele.appendChild(aside);
    }
    data.length ? (ele.style.display = "block") : null;
    document.querySelector("section").innerText = `Found ${data.length} results in ${timeEnd - timeStart}ms using trie`;
};