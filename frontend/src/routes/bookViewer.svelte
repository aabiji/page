<script lang="ts">
    import { onMount } from "svelte";
    import { EpubViewer } from "./epub";
    import { callApi } from "./utils";

    function getBookInfo(name: string) {
        return callApi(`http://localhost:8080/book/get/${name}`, "GET", {});
    }

    let errorOut = false;
    let book = {
        TableOfContents: [],
        Info: {
            Title: "-", Author: "-", Description: "-", Date: "-",
            Contributor: "-", Coverage: "-", Source: "-", Rights: "-",
            Relation: "-", Publisher: "0", Language: "-", "Identifier": "-",
            Subjects: []
        },
    };
    onMount(() => {
        let div = document.getElementById("book-view")!;

        getBookInfo("AnimalFarm").then((json) => {
            if ("Server error" in json) {
                errorOut = true;
                return;
            }

            book = json;
            book.Info.Subjects = book.Info.Subjects == null ? [] : book.Info.Subjects;
            let e = new EpubViewer(json.Files, div);
            e.render();
        });
    });
</script>

{#if errorOut}
<div class="error">
    <p> Oops, something went wrong </p>
</div>
{:else}
<div class="container">
    <div class="left-sidepanel">
        <h1> {book.Info.Title} </h1>
        <!--- cover image goes here: --->
        <h3> {book.Info.Author} </h3>
        <h5> {book.Info.Description} </h5>
        <p> Date: {book.Info.Date} </p>
        <p> Contributor: {book.Info.Contributor} </p>
        <p> Coverage: {book.Info.Coverage} </p>
        <p> Source: {book.Info.Source} </p>
        <p> Rights: {book.Info.Rights} </p>
        <p> Relation: {book.Info.Relation} </p>
        <p> Publisher: {book.Info.Publisher} </p>
        <p> Language: {book.Info.Language} </p>
        <p> Identifier: {book.Info.Identifier} </p>
        <p> Subjects: {#each book.Info.Subjects as subject} {subject}, {/each} </p>
        <hr>
        <h3> Table of contents </h3>
        <ol>
            {#each book.TableOfContents as section}
                <li><a href={section[1]}>{section[0]}</a></li>
            {/each}
        </ol>
    </div>
    <div class="right-sidepanel">
        <div id="book-view"></div>
    </div>
</div>
{/if}

<style>
    .container {
        display: flex;
    }

    #book-view {
        width: 55%;
        margin: 0 auto;
    }

    .right-sidepanel {
        width: 75%;
        border: 1px solid black;
    }

    .left-sidepanel {
        width: 25%;
        height: 100vh;
        overflow-y: scroll;
        background-color: #a8a8a8;
        overflow-wrap: break-word;
    }

    .left-sidepanel h1, h3, h5 {
        text-align: center;
    }
</style>
