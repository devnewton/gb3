class Gb3Coincoin extends HTMLElement {

    constructor() {
        super();
    }

    connectedCallback() {
        let postsElement = document.createElement('gb3-posts');
        this.appendChild(postsElement);

        var postSource = new EventSource("/api/poll");
        postSource.onmessage = function (event) {
            console.log(event.data);
            let post = JSON.parse(event.data);
            postsElement.insertPost(post);
        };
    }

}
customElements.define('gb3-coincoin', Gb3Coincoin);