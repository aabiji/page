<script lang="ts">
    import { onMount } from "svelte";
    import Navbar from "./upload.svelte";

    let showUploadDialog: boolean;
    let userProfilePic: string;
    onMount(() => {
        userProfilePic = `${window.origin}/default-profile.png`;

        document.onclick = (event) => {
            if (event.target == null) return;
            let t = event.target as HTMLElement;
            let dialog = document.getElementsByTagName("dialog")[0];
            if (dialog == undefined) return;
            if (!t.className.includes("upload") && t != dialog && !dialog.contains(t)) {
                showUploadDialog = false;
            }
        }
    });
</script>

<div class="bar">
    <h1 title="Home" class="logo"><a href="/">Page</a></h1> 
    <input title="Search" type="text" placeholder="Search books">
    <div class="actions">
        <button title="Upload book" class="upload"
         on:click={() => showUploadDialog = !showUploadDialog}> + </button>
        <a title="Account settings" href="/account"><img alt="Profile" src={userProfilePic} /></a>
    </div>
</div>
{#if showUploadDialog}
<div>
    <Navbar />
</div>
{/if}

<style>
    .bar {
        top: 0;
        height: var(--navbar-height);
        width: 100.5%;
        position: fixed;
        padding: 5px 5px;
        text-align: center;
        margin-left: -10px;
        background-color: var(--background-accent);
    }

    .bar * {
        display: inline;
    }

    .actions {
        float: right;
    }

    h1 {
        float: left;
        margin-top: -1px;
        margin-left: 10px;
        margin-bottom: 0px;
    }

    button {
        font-size: 25px;
        color: var(--accent-color);
        background-color: rgba(0,0,0,0);
        border: 1px solid #535454;
        margin-top: 10px !important;
        border-radius: 25px;
        padding: 3px 10px;
        top: -8px;
        position: relative;
    }

    button:hover {
        background-color: #535454;
        border: 1px solid var(--background-accent);
    }

    img {
        width: 35px;
        height: 35px;
        margin-top: 5px;
        margin-left: 20px;
        margin-right: 20px;
        border-radius: 50%;
    }

    input {
        width: 50%;
        color: white;
        font-size: 18px;
        padding: 8px 8px;
        border: #535454 1px solid;
        background-color: rgba(0,0,0,0);
        border-radius: 10px;
    }
 
    input:hover, input:focus {
        outline: none;
        border: var(--accent-color) 1px solid;
    }
  
    input:-webkit-autofill, input:-webkit-autofill:hover, input:-webkit-autofill:focus {
        font-size: 18px;
        -webkit-text-fill-color: white;
        -webkit-box-shadow: 0 0 0px 40rem #1b1c1c inset;
    }
</style>