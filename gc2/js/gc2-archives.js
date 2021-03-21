class Gc2Archives extends HTMLElement {

    constructor() {
        super();
        this.tribunes = new Map();
    }

    async connectedCallback() {
        this.setupControls();
        this.setupResultsContainer();
        await this.setupBackend2html();
    }


    async setupBackend2html() {
        let response = await fetch("/peg/backend2html.pegjs");
        let text = await response.text();
        this.backend2html = peg.generate(text);
    }

    setupControls() {
        let form = document.createElement("form");
        this.appendChild(form);

        let termsInput = document.createElement("input");
        termsInput.type = "text";
        termsInput.spellcheck = true;
        form.appendChild(termsInput);

        let searchButton = document.createElement('button');
        searchButton.innerText = "Search";
        searchButton.type = "submit";
        form.appendChild(searchButton);

        let backButton = document.createElement("button");
        backButton.innerText = "Back"
        backButton.onclick = () => {
            window.location.href = "/";
        }
        form.appendChild(backButton);

        form.onsubmit = (e) => {
            e.preventDefault()
            if (termsInput.value) {
                termsInput.classList.toggle("gc2-loading", true);
                fetch(`/gb2c/search?terms=${encodeURIComponent(termsInput.value)}`, {
                    method: "GET",
                }).then((response) => {
                    termsInput.classList.toggle("gc2-loading", false);
                    return response.json();
                }).then((data) => {
                    this.setResults(data);
                }).catch((error) => {
                    termsInput.classList.toggle("gc2-loading", false);
                    console.log(`Cannot search archives. Error: `, error);
                });
            }
        };
    }

    setupResultsContainer() {
        this.resultsContainer = document.createElement("gc2-tribune");
        this.resultsContainer.classList.add("gc2-archives-search-results");
        this.appendChild(this.resultsContainer);
    }

    setResults(results) {
        this.clearResults();
        for (let post of results.posts) {
            this.insertPost(post);
        }
    }

    insertPost(post) {
        if (!document.getElementById(`${post.id}@${post.tribune}`)) {
            let postElement = new Gc2Post();
            post.message = this.backend2html.parse(post.message);
            postElement.setupWithTribune(post);
            this.resultsContainer.appendChild(postElement);
        }
    }

    clearResults() {
        let child;
        while (child = this.resultsContainer.firstChild) {
            this.resultsContainer.removeChild(child);
        }
    }

}
customElements.define('gc2-archives', Gc2Archives);
