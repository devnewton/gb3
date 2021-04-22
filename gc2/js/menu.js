function gc2CloseMenu() {
    let menu = document.getElementsByTagName("gc2-menu")[0];
    menu.style.display = "none";
    menu.clear();
    document.querySelector("gc2-main").style.display = "";
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
        totozButton.innerText = "😱 Totoz";
        totozButton.onclick = () => {
            this.showComponent("gc2-totozsearch");
        }
        this.appendChild(totozButton);

        let emojiButton = document.createElement('button');
        emojiButton.innerText = "💩 Emoji";
        emojiButton.onclick = () => {
            this.showComponent("gc2-emojisearch");
        }
        this.appendChild(emojiButton);

        let attachButton = document.createElement('button');
        attachButton.innerText = "📎 Attach";
        attachButton.onclick = () => {
            this.showComponent("gc2-attach");
        }
        this.appendChild(attachButton);

        let settingsButton = document.createElement('button');
        settingsButton.innerText = "⚙ Settings";
        settingsButton.onclick = () => {
            window.location.href = "/settings.html";
        }
        this.appendChild(settingsButton);

        let archivesButton = document.createElement('button');
        archivesButton.innerText = "🔍 Archives";
        archivesButton.onclick = () => {
            window.location.href = "/archives.html";
        }
        this.appendChild(archivesButton);

        let backButton = document.createElement('button');
        backButton.innerText = "↩ Back";
        backButton.onclick = () => {
            gc2CloseMenu();
        }
        this.appendChild(backButton);
    }

    showComponent(componentName) {
        this.clear();

        let component = document.createElement(componentName);
        component.setup();
        this.appendChild(component);

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
        let form = document.createElement("form");

        let searchInput = document.createElement('input');
        searchInput.type = "text";
        searchInput.placeholder = "dont be so vanilla";
        form.appendChild(searchInput);

        let searchButton = document.createElement('button');
        searchButton.innerText = "Search";
        searchButton.type = "submit";
        form.appendChild(searchButton);

        form.onsubmit = (e) => {
            e.preventDefault();
            fetch(`/gb2c/totoz/search?terms=${encodeURIComponent(searchInput.value)}`, {
                method: "GET",
            }).then((response) => {
                return response.json();
            }).then((data) => {
                this.setResults(data);
            }).catch((error) => {
                console.log(`Cannot search totoz. Error: `, error);
            });
        };
        this.appendChild(form);


        this.resultsContainer = document.createElement("div");
        this.resultsContainer.onclick = (e) => {
            let caption = e.target.parentElement.querySelector("figcaption");
            if (caption) {
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
        for (let totoz of results.totozes) {
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

class Gc2EmojiSearch extends HTMLElement {
    constructor() {
        super();
    }

    setup() {
        let form = document.createElement("form");

        let searchInput = document.createElement('input');
        searchInput.type = "text";
        searchInput.placeholder = "poop";
        form.appendChild(searchInput);

        let searchButton = document.createElement('button');
        searchButton.innerText = "Search";
        form.appendChild(searchButton);

        form.onsubmit = (e) => {
            e.preventDefault();
            fetch(`/gb2c/emoji/search?terms=${encodeURIComponent(searchInput.value)}`, {
                method: "GET",
            }).then((response) => {
                return response.json();
            }).then((data) => {
                this.setResults(data);
            }).catch((error) => {
                console.log(`Cannot search emoji. Error: `, error);
            });
        };
        this.appendChild(form);

        this.resultsContainer = document.createElement("div");
        this.resultsContainer.onclick = (e) => {
            let characters = e.target.parentElement.querySelector(".gc2-emoji-characters");
            if (characters) {
                let message = document.getElementById("gc2-message");
                message.value += `${message.value && ' '}${characters.innerText} `;
                gc2CloseMenu();
                message.focus();
            }
        };
        this.resultsContainer.classList.add("gc2-emoji-search-results")
        this.appendChild(this.resultsContainer);
    }

    setResults(results) {
        this.clearResults();
        for (let emoji of results) {
            let emojiElement = document.createElement("figure");

            let emojiCharacters = document.createElement("p");
            emojiCharacters.innerText = emoji.characters;
            emojiCharacters.classList.add("gc2-emoji-characters")
            emojiElement.appendChild(emojiCharacters);

            let emojiName = document.createElement("figcaption");
            emojiName.innerText = emoji.name;
            emojiElement.appendChild(emojiName);

            this.resultsContainer.appendChild(emojiElement);
        }
    }

    clearResults() {
        let child;
        while (child = this.resultsContainer.firstChild) {
            this.resultsContainer.removeChild(child);
        }
    }
}
customElements.define("gc2-emojisearch", Gc2EmojiSearch);

class Gc2Attach extends HTMLElement {
    constructor() {
        super();
    }

    setup() {
        let form = document.createElement("form");
        form.enctype = "multipart/form-data";
        form.method = "post";

        let fileInput = document.createElement('input');
        fileInput.name = "attachment";
        fileInput.type = "file";
        form.appendChild(fileInput);

        let attachButton = document.createElement('button');
        attachButton.innerText = "Attach";
        form.appendChild(attachButton);

        let progress = document.createElement("img");
        progress.classList.add("gc2-attach-progress")

        form.onsubmit = (e) => {
            e.preventDefault();
            progress.src = "/img/loading.svg";
            fetch(`/gb2c/attachment/`, {
                method: "POST",
                body: new FormData(form)
            }).then((response) => {
                progress.src = "";
                let location = response.headers.get("Location");
                if (location) {
                    let message = document.getElementById("gc2-message");
                    message.value += `${message.value && ' '}${location} `;
                    gc2CloseMenu();
                    message.focus();
                }
            }).catch((error) => {
                progress.src = "";
                console.log(`Cannot attach file. Error: `, error);
            });
        };
        this.appendChild(form);

        this.appendChild(progress);

        let imagePreview = document.createElement("img");
        this.appendChild(imagePreview);
        let audioPreview = document.createElement("audio");
        audioPreview.controls = true;
        this.appendChild(audioPreview);
        let videoPreview = document.createElement("video");
        videoPreview.controls = true;
        this.appendChild(videoPreview);

        let updatePreview = () => {
            imagePreview.style.display = "none";
            audioPreview.style.display = "none";
            audioPreview.pause();
            audioPreview.removeAttribute("src");
            audioPreview.load();
            videoPreview.style.display = "none";
            videoPreview.pause();
            videoPreview.removeAttribute("src");
            videoPreview.load();
            if (fileInput.files && fileInput.files.length === 1) {
                let file = fileInput.files[0];
                let type = file.type;
                if (type.startsWith("image")) {
                    imagePreview.src = window.URL.createObjectURL(file);
                    imagePreview.style.display = "";
                } else if (type.startsWith("audio")) {
                    audioPreview.src = window.URL.createObjectURL(file);
                    audioPreview.style.display = "";
                } else if (type.startsWith("video")) {
                    videoPreview.src = window.URL.createObjectURL(file);
                    videoPreview.style.display = "";
                }
            }
        };

        updatePreview();
        fileInput.onchange = () => updatePreview();

        document.addEventListener('paste', (event) => {
            if (!(form.offsetWidth || form.offsetHeight || form.getClientRects().length)) {
                return;
            }
            let files = (event.clipboardData || event.originalEvent.clipboardData).files;
            if (files && files.length === 1) {
                fileInput.files = files;
                updatePreview();
            }
            event.preventDefault();
        });
    }
}
customElements.define("gc2-attach", Gc2Attach);