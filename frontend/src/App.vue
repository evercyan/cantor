<template>
  <div id="app" @contextmenu.prevent>
    <Cantor />
  </div>
</template>

<script>
import "@/assets/css/main.less";
import Cantor from "@/components/Cantor.vue";
export default {
  name: "app",
  components: {
    Cantor,
  },
};
</script>