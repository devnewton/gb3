function gc2CloseMenu() {
    let menu = document.getElementsByTagName("gc2-menu")[0];
    menu.style.display = "none";
    menu.clear();
    document.querySelector("gc2-main").style.display = "flex";
}

class Gc2Menu extends HTMLElement {
    constructor() {
        super();
        this.style.display = "none";
        this.showSelector();
    }

    showSelector() {
        this.clear();

        let totozButton = document.createElement('button');
        totozButton.innerText = "Totoz";
        totozButton.onclick = () => {
            this.showTotoz();
        }
        this.appendChild(totozButton);

        let emojiButton = document.createElement('button');
        emojiButton.innerText = "Emoji";
        this.appendChild(emojiButton);

        let binButton = document.createElement('button');
        binButton.innerText = "Bin";
        this.appendChild(binButton);

        let settingsLink = document.createElement('a');
        settingsLink.innerText = "Settings";
        settingsLink.href = "/settings.html";
        this.appendChild(settingsLink);

        let closeButton = document.createElement('button');
        closeButton.innerText = "Close";
        closeButton.onclick = () => {
            gc2CloseMenu();
        }
        this.appendChild(closeButton);
    }

    showTotoz() {
        this.clear();

        let totozSearch = document.createElement("gc2-totozsearch");
        totozSearch.setup();
        this.appendChild(totozSearch);

        let backButton = document.createElement("button");
        backButton.innerText = "Back";
        backButton.onclick = () => {
            this.showSelector();
        }
        this.appendChild(backButton);
    }

    clear() {
        let child;
        while (child = this.firstChild) {
            this.removeChild(child);
        }
    }
}
customElements.define('gc2-menu', Gc2Menu);

class Gc2TotozSearch extends HTMLElement {
    constructor() {
        super();
    }

    setup() {
        let searchInput = document.createElement('input');
        searchInput.type = "text";
        searchInput.placeholder = "dont be so vanilla";
        this.appendChild(searchInput);

        let searchButton = document.createElement('button');
        searchButton.innerText = "Search";
        searchButton.onclick = () => {
            fetch(`/api/totoz/search?terms=${encodeURIComponent(searchInput.value)}`, {
                method: "GET",
            }).then((response) =>{
                return response.json();
            }).then((data) =>{
                this.setResults(data);
            }).catch((error) => {
                console.log(`Cannot search totoz. Error: `, error);
            });
        }
        this.appendChild(searchButton);

        this.resultsContainer = document.createElement("div");
        this.resultsContainer.onclick = (e) => {
            let caption = e.target.parentElement.querySelector("figcaption");
            if(caption) {
                let message = document.getElementById("gc2-message");
                let value = message.value;
                message.value += `${message.value && ' '}[:${caption.innerText}] `;
                gc2CloseMenu();
                message.focus();
            }
        };
        this.resultsContainer.classList.add("gc2-totoz-search-results")
        this.appendChild(this.resultsContainer);
    }

    setResults(results) {
        this.clearResults();
        for(let totoz of results.totozes) {
            let totozElement = document.createElement("figure");

            let totozImg = document.createElement("img");
            totozImg.src = totoz.image;
            totozElement.appendChild(totozImg);

            let totozName = document.createElement("figcaption");
            totozName.innerText = totoz.name;
            totozElement.appendChild(totozName);

            this.resultsContainer.appendChild(totozElement);
        }
    }

    clearResults() {
        let child;
        while (child = this.resultsContainer.firstChild) {
            this.resultsContainer.removeChild(child);
        }
    }
}
customElements.define("gc2-totozsearch", Gc2TotozSearch);