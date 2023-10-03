<script lang="ts">
    import BookViewer from "../bookview.svelte";
    import Navbar from "../../../components/navbar.svelte";
    import * as utils from "$lib/utils";

    import { onMount } from "svelte";
    import { page } from "$app/stores";

    onMount(() => {
        utils.redirectIfNotAuth();
    });
</script>

<Navbar />
<BookViewer bookId={parseInt($page.params.bookId)} />