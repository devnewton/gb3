class Gb3Posts extends HTMLElement {

    constructor() {
        super();
    }

    insertPost(post) {
        let postElements = this.getElementsByTagName('gb3-post');
        for (let e of postElements) {
            if (e.tribune === post.tribune && post.id > e.postId) {
                this.insertPostAfterElement(post, e);
                return;
            }
        }
        let postElement = new Gb3Post();
        postElement.setup(post);
        this.appendChild(postElement);
    }

    insertPostAfterElement(post, element) {
        let postElement = new Gb3Post();
        postElement.setup(post);
        this.insertBefore(postElement, element.nextSibling);
    }

}
customElements.define('gb3-posts', Gb3Posts);