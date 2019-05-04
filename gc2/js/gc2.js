class Gc2Main extends HTMLElement {

    constructor() {
        super();
        this.tribunes = new Set();
    }

    connectedCallback() {
        this.setupTribuneSelect();
        this.setupMessageInput();
        this.setupPostsElement();
        this.setupBackend2html();
    }

    setupTribuneSelect() {
        this.tribuneSelect = document.createElement("select");
        let viewAllTribuneOption = document.createElement("option");
        viewAllTribuneOption.value = "";
        viewAllTribuneOption.text = "All";
        this.tribuneSelect.add(viewAllTribuneOption);
        document.styleSheets[0].insertRule("")
        this.tribuneSelect.onchange = () => {
            this.postsElement.dataset.tribune = this.tribuneSelect.value;
        }
        this.appendChild(this.tribuneSelect);
    }

    setupMessageInput() {
        this.messageInput = document.createElement("input");
        this.messageInput.type = "text";
        this.messageInput.spellcheck = true;
        this.appendChild(this.messageInput);
    }

    setupPostsElement() {
        this.postsElement = document.createElement('gc2-posts');
        this.appendChild(this.postsElement);
    }

    setupBackend2html() {
        fetch("/peg/backend2html.pegjs")
            .then((response) => {
                return response.text();
            })
            .then((text) => {
                this.backend2html = peg.generate(text);
                this.startPoll();
            });
    }

    startPoll() {
        let postSource = new EventSource("/api/poll");
        postSource.onmessage = (event) => {
            let post = JSON.parse(event.data);
            post.message = this.backend2html.parse(post.message);
            this.postsElement.insertPost(post);
            this.addTribune(post.tribune);
        };
    }

    addTribune(tribune) {
        if (!this.tribunes.has(tribune)) {
            this.tribunes.add(tribune);
            var option = document.createElement("option");
            option.text = option.value = tribune;
            this.tribuneSelect.add(option);
        }
    }

}
customElements.define('gc2-main', Gc2Main);