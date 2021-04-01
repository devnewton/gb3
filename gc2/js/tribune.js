class Gc2Tribune extends HTMLElement {

    constructor() {
        super();
    }

    insertPost(post) {
        if (!document.getElementById(`${post.id}@${post.tribune}`)) {
            let postElement = new Gc2Post();
            postElement.setup(post);
            this.appendChild(postElement);
            this.sortPosts();
        }
    }

    sortPosts() {
        let postElements = this.getElementsByTagName('gc2-post');
        Array.prototype.slice.call(postElements).sort((a, b) => {
            if (a.postId < b.postId) {
                return -1;
            } else if (a.postId > b.postId) {
                return 1;
            }
            return 0;
        }).forEach((p) => {
            this.appendChild(p);
        });
    }
}
customElements.define('gc2-tribune', Gc2Tribune);

class Gc2Tribunes extends HTMLElement {
    constructor() {
        super();
    }
}
customElements.define('gc2-tribunes', Gc2Tribunes);