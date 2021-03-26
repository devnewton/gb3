class Gc2Tribune extends HTMLElement {

    constructor() {
        super();
        let posterOrder = localStorage.postOrder || "chronological";
        this.chronologicalOrder = posterOrder === "chronological";
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
            if(this.chronologicalOrder) {
                this.appendChild(p);
            } else {
                this.prepend(p);
            }
        });
    }
}
customElements.define('gc2-tribune', Gc2Tribune);