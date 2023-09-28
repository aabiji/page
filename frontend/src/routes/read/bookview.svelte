<script lang="ts">
    import { onMount } from "svelte";
    import { writable } from "svelte/store";

    import { EpubViewer } from "./epub";
    import * as utils from "$lib/utils";

    export let bookName: string;

    let errorOut = false;
    let bookView: HTMLElement;
    let toggleButton: HTMLElement;
    let leftSidepanl: HTMLElement;
    let book = writable({
        CurrentPage: 0,
        Epub: {
            CoverImagePath: "",
            Info: {
                Author: "",
                Contributor: "",
                Coverage: "",
                Date: "",
                Description: "",
                Identifier: "",
                Language: "",
                Publisher: "",
                Relation: "",
                Rights: "",
                Source: "",
                Subjects: [],
                Title: "",
            },
            TableOfContents: [{
                Path: "",
                Name: "",
            }],
        },
    });
    let epub = writable(new EpubViewer([], [], 0));

    function toggelLeftSidepanel() {
        toggleButton.classList.toggle("left");
        leftSidepanl.classList.toggle("hidden-left-sidepanel");
    }

    onMount(() => {
        utils.callApi(`http://localhost:8080/book/get/${bookName}`, "GET", {}).then((response) => {
            if (utils.serverError in response) {
                errorOut = true;
                console.log(response[utils.serverError]);
                return;
            }

            response.Epub.CoverImagePath = utils.coverImagePath(response.Epub.CoverImagePath);
            book.set(response);

            let e = new EpubViewer(response.FileScrollOffsets, response.Epub.Files, response.CurrentPage, bookView);
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
        <h1> {$book.Epub.Info.Title} </h1>
        <img alt="Ebook cover image" src={$book.Epub.CoverImagePath}/>
        <h3> {$book.Epub.Info.Author} </h3>
        <hr>
        <h5> {$book.Epub.Info.Description} </h5>
        <p> Date: {$book.Epub.Info.Date} </p>
        <p> Contributor: {$book.Epub.Info.Contributor} </p>
        <p> Coverage: {$book.Epub.Info.Coverage} </p>
        <p> Source: {$book.Epub.Info.Source} </p>
        <p> Rights: {$book.Epub.Info.Rights} </p>
        <p> Relation: {$book.Epub.Info.Relation} </p>
        <p> Publisher: {$book.Epub.Info.Publisher} </p>
        <p> Language: {$book.Epub.Info.Language} </p>
        <p> Identifier: {$book.Epub.Info.Identifier} </p>
        <p> Subjects: {#each $book.Epub.Info.Subjects as subject} {subject}  {/each} </p>
        <hr>
        <h3> Table of contents </h3>
        <ol>
            {#each $book.Epub.TableOfContents as section}
                <li><span on:click={$epub.jumpToSection(section.Path)}>{section.Name}</span></li>
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
    span {
        cursor: pointer;
        text-decoration: none;
        color: var(--accent-color);
    }
    span:hover {
        color: var(--accent-color-darken);
    }

    .container {
        width: 100%;
        position: fixed;
        margin-top: 45px;
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
