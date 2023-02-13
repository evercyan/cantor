<template>
  <div id="app" @contextmenu.prevent style="--wails-draggable:drag">
    <Cantor/>
  </div>
</template>

<script>
import "@/assets/css/main.less";
import Cantor from "@/components/cantor.vue";

export default {
  name: "app",
  components: {
    Cantor,
  },
}
</script>