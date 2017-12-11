var vm = new Vue({
  el: '#app',
  template: '<div><div class="counter">{{ counterVal }}</div><button class="btn btn-primary" v-on:click="increment">Increment</button><div>{{ incrementText }}</div></div>',
  data: {
    counterVal: 0,
    incrementText: ""
  },
  created: function() {
    var vm = this;
    goui.onMessage("add", function(message){
      vm.counterVal = parseInt(message);
    });

    goui.onMessage("getIncText", function(message){
      vm.incrementText = message;
    });
  },
  methods: {
    increment: function() {      
      var prevVal = this.counterVal 
      goui.send("add", {
        val1: this.counterVal, 
        val2: 1
      });
      
      goui.send("getIncText", {
        val1: prevVal, 
        val2:1
      });      
    },
  }
});
