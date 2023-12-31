<script lang="ts">
    import { onMount } from "svelte";
    import { writable } from "svelte/store";

    import { EpubViewer } from "./epub";
    import * as utils from "$lib/utils";

    export let bookId: number;

    let errorOut = false;
    let bookView: HTMLElement;
    let toggleButton: HTMLElement;
    let leftSidepanl: HTMLElement;
    let book = writable(new utils.Book());
    let epub = writable(new EpubViewer([], [], 0));

    function toggelLeftSidepanel() {
        toggleButton.classList.toggle("left");
        leftSidepanl.classList.toggle("hidden-left-sidepanel");
    }

    onMount(() => {
        let bookJson = utils.cacheGet(utils.BookKey(bookId));
        book.set(bookJson);

        let url = `${utils.backendOrigin}/user/book/get/${bookId}`;
        utils.callApi(url, "GET").then((response) => {
            let e = new EpubViewer(response.ScrollOffsets, bookJson.Files, response.CurrentPage, bookView)
            epub.set(e);
            $epub.render();
        });
    });
</script>

{#if errorOut}
<div class="error">
    <p> Oops, something went wrong </p>
</div>
{:else}
<div class="container">
    <div class="left-sidepanel" bind:this={leftSidepanl}>
        <h1> {$book.Info.Title} </h1>
        <img alt="Ebook cover" src={$book.CoverImagePath}/>
        <h3> {$book.Info.Author} </h3>
        <hr>
        <h5> {$book.Info.Description} </h5>
        <p> Date: {$book.Info.Date} </p>
        <p> Contributor: {$book.Info.Contributor} </p>
        <p> Coverage: {$book.Info.Coverage} </p>
        <p> Source: {$book.Info.Source} </p>
        <p> Rights: {$book.Info.Rights} </p>
        <p> Relation: {$book.Info.Relation} </p>
        <p> Publisher: {$book.Info.Publisher} </p>
        <p> Language: {$book.Info.Language} </p>
        <p> Identifier: {$book.Info.Identifier} </p>
        <p> Subjects: {#each $book.Info.Subjects as subject} {subject}  {/each} </p>
        <hr>
        <h3> Table of contents </h3>
        <ol>
            {#each $book.TableOfContents as section}
                <li><button on:click={() => $epub.jumpToSection(section.Path)}>
                    {section.Name}
                </button></li>
            {/each}
        </ol>
    </div>
    <button class="left-sidepanel-toggle"
            bind:this={toggleButton}
            title="Toggle sidepanel visiblity"
            on:click={toggelLeftSidepanel}>
        &gt;
    </button>
    <div class="right-sidepanel">
        <div bind:this={bookView} id="book-view"></div>
    </div>
</div>
{/if}

<style>
    button {
        cursor: pointer;
        text-decoration: none;
        color: var(--accent-color);
        background-color: var(--background-accent);
        font-size: 16px;
    }
 
    button:hover {
        color: var(--accent-color-darken);
    }

    .container {
        width: 100%;
        position: fixed;
        margin-top: var(--navbar-height);
        display: flex;
        align-items: center;
        justify-content: center;
    }

    #book-view {
        width: 600px;
        height: 100%;
        margin: 0 auto;
    }

    .right-sidepanel {
        width: 80%;
        height: calc(97.5vh - 45px);
    }

    .left-sidepanel {
        width: 20%;
        margin-left: -5px;
        overflow-y: scroll;
        height: calc(98vh - 45px);
        overflow-wrap: break-word;
        background-color: var(--background-accent);
    }

    .left-sidepanel p {
        font-size: 15px;
    }

    .left-sidepanel img {
        width: auto;
        height: 225px;
        margin-left: 25%;
    }

    .left-sidepanel h1, h3, h5 {
        text-align: center;
    }

    .left-sidepanel-toggle {
        height: 35px;
        border: none;
        color: white;
        cursor: pointer;
        margin-left: -10px;
        align-content: flex-start;
        background-color: var(--accent-color);
    }

    .left-sidepanel-toggle:hover {
        background-color: var(--accent-color-darken);
    }

    :global(.left-sidepanel-toggle.left) {
        left: 0px;
        position: absolute;
        margin-left: 0px !important;
    }

    :global(.hidden-left-sidepanel) {
        display: none;
    }
</style>