<script lang="ts">
    import { onMount } from "svelte";
    import { writable } from "svelte/store";

    import { EpubViewer } from "./epub";
    import * as utils from "./utils";

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

    function jumpToSection(sectionPath: string) {
        let files = $epub.files.map((f: File) => f.Path);
        let sectionParts = sectionPath.split("#");
        let section = sectionParts[1] == undefined ? "" : sectionParts[1];
        let index = files.indexOf(sectionParts[0]);
        $epub.pageIdx = index;
        $epub.sections[index] = section;
        $epub.render();
    }

    function getBook(name: string) {
        utils.callApi(`http://localhost:8080/book/get/${name}`, "GET", {}).then((json) => {
            if ("Server error" in json) {
                errorOut = true;
                console.log(json);
                return;
            }

            if (json.Epub.CoverImagePath == "") {
                json.Epub.CoverImagePath = "default-cover-image.png";
            } else {
                json.Epub.CoverImagePath = utils.staticFileUrl(json.Epub.CoverImagePath);
            }
            book.set(json);

            epub.set(new EpubViewer(json.FileScrollOffsets, json.Epub.Files, json.CurrentPage, bookView));
            $epub.render();
        });
    }

    onMount(() => {
        utils.callApi("http://localhost:8080/cookie", "GET", {}).then((() => {
            getBook("med");
        }));
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
                <li><span on:click={jumpToSection(section.Path)}>{section.Name}</span></li>
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
        color: #4287f5;
        cursor: pointer;
        text-decoration: none;
    }
    span:hover {
        color: #5a98fa;
    }

    .container {
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
        height: 97vh;
    }

    .left-sidepanel {
        width: 20%;
        height: 97vh;
        overflow-y: scroll;
        background-color: #1c1c1c;
        overflow-wrap: break-word;
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
        background-color: #1e63d4;
    }

    .left-sidepanel-toggle:hover {
        background-color: #1757bf;
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
