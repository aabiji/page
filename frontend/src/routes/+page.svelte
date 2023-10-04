<script lang="ts">
    import { onMount } from "svelte";
    import { writable } from "svelte/store";
 
    import Book from "../components/book.svelte";
    import Navbar from "../components/navbar.svelte";
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

    function removeBook(id: number) {
        let url = `${utils.backendOrigin}/user/book/remove/${id}`;
        utils.callApi(url, "POST").then((response) => {
            if (utils.serverError in response) {
                console.log(response);
                return;
            }
            utils.removeBook(id);
            loadBooks(utils.cacheGet(utils.BooksKey));
        });
    }

    onMount(() => {
        utils.redirectIfNotAuth();
        loadBooks(utils.cacheGet(utils.BooksKey));
    });
</script>

<Navbar />
<div class="container">
    <h1> Your books </h1>
    {#if $books.length == 0}
        <p> Looks like you don't have any books yet! </p>
    {/if}
    <div class="collection">
        {#each $books as b}
            <div class="book">
                <Book cover={b.cover} title={b.title} id={b.id} />
                <button on:click={() => removeBook(b.id)}> Remove book </button>
            </div>
        {/each}
    </div>
</div>

<style>
    .book {
        text-align: center;
    }

    .book button {
        padding: 8px 8px;
        color: white;
        border: none;
        font-size: 14px;
        background-color: red;
        margin-bottom: 15px;
        display: none;
        margin-left: 12px;
    }

    .book button:hover {
        background-color: darkred;
    }

    .book:hover button {
        display: block;
    }

    .container {
        padding: 10px;
        margin-top: var(--navbar-height);
    }

    .collection {
        gap: 30px;
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(115px, 1fr));
    }
</style>