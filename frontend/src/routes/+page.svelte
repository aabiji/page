<script lang="ts">
    import { onMount } from "svelte";
    import { writable } from "svelte/store";
 
    import Book from "./book.svelte";
    import Navbar from "./navbar.svelte";
    import * as utils from "$lib/utils";
 
    interface BookDisplayInfo {
        id: number,
        cover: string,
        title: string,
    };
 
    let books = writable<BookDisplayInfo[]>([]);
    function loadBooks(bookIds: number[]) {
        $books = [];
        for (let i = 0; i < bookIds.length; i++) {
            let id = bookIds[i];
            let obj: utils.Book = utils.cacheGet(utils.BookKey(id));
            let display: BookDisplayInfo = {
                id: id,
                title: obj.Info.Title,
                cover: obj.CoverImagePath,    
            };
            $books.push(display);
        }
    }

    function addBook(id: number) {
        let url = `${utils.backendOrigin}/book/get/${id}`
        utils.callApi(url, "GET").then((info: utils.Book) => {
            info.CoverImagePath = utils.coverImagePath(info.CoverImagePath);
            utils.cacheBook(id, info);
            loadBooks(utils.cacheGet(utils.BooksKey));
        });
    }

    let fileInput: HTMLElement;
    function uploadFile(event: any) {
        const file = event.target.files[0];
        const formData = new FormData();
        formData.append("file", file);
        let url = `${utils.backendOrigin}/user/book/upload`;
        utils.callApi(url, "POST", formData, true).then((response) => {
            if (utils.serverError in response) return;
            addBook(response.BookId);
        });
    }

    // Drap and drop files
    const disableDrag = (event: DragEvent) => event.preventDefault();
    function dropFile(event: DragEvent) {
        event.preventDefault();
        let files = event.dataTransfer?.files;
        if (files == undefined) return;
        for (let i = 0; i < files.length; i++) {
            let file = files[i];
            console.log(file);
        }
    }

    onMount(() => {
        utils.redirectIfNotAuth();
        loadBooks(utils.cacheGet(utils.BooksKey));
    });
</script>

<Navbar />

<dialog open>
    <h1> Upload a book </h1>
    <input on:change={uploadFile} bind:this={fileInput} type="file" style="display: none;">
    <button on:click={() => fileInput.click()}> Upload book </button>
    <div class="drop" on:drop={dropFile} on:dragover={disableDrag}>
        <p> drop a file here </p>
    </div>
</dialog>

<div class="container">
    <div class="top">
        <h1> Your books </h1>
        <input on:change={uploadFile} bind:this={fileInput} type="file" style="display: none;">
        <button on:click={() => fileInput.click()}> Upload book </button>
    </div>
    {#if $books.length == 0}
        <p> Looks like you don't have any books yet! </p>
    {/if}
    <div class="collection">
        {#each $books as b}
            <Book cover={b.cover} title={b.title} id={b.id} />
        {/each}
    </div>
</div>

<style>
    button {
        color: white;
        margin-left: 10px;
        font-size: 18px;
        padding: 5px 5px;
        margin-bottom: 15px;
        background-color: var(--accent-color);
    }
    button:hover {
        background-color: var(--accent-color-darken);
    }
    .container {
        padding: 10px;
        margin-top: var(--navbar-height);
    }
    .top * {
        display: inline;
    }
    .top {
        margin-bottom: 20px;
    }
    .collection {
        gap: 30px;
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(115px, 1fr));
    }

    dialog {
        width: 300px;
        height: 300px;
        text-align: center;
        position: absolute;
    }

    .drop {
        height: 200px;
        background-color: red;
        border: 3px dashed black;
    }
</style>