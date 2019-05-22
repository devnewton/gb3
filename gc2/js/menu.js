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
            this.style.display = "none";
            document.querySelector("gc2-main").style.display = "flex";
        }
        this.appendChild(closeButton);
    }

    showTotoz() {
        this.clear();

        let searchInput = document.createElement('input');
        searchInput.type = "text";
        searchInput.placeholder = "dont be so vanilla";
        this.appendChild(searchInput);

        let searchButton = document.createElement('button');
        searchButton.innerText = "Search";
        searchButton.onclick = () => {
            //TODO
        }
        this.appendChild(searchButton);

        let backButton = document.createElement('button');
        backButton.innerText = "Back";
        backButton.onclick = () => {
            this.showSelector();
        }
        this.appendChild(backButton);
    }

    clear() {
        let child;
        while(child = this.firstChild) {
            this.removeChild(child);
        }
    }
}
customElements.define('gc2-menu', Gc2Menu);