var vm = new Vue({
  el: '#app',
  template: '<div><div class="counter">{{ counterVal }}</div><button class="btn btn-primary" v-on:click="increment">Increment</button><div>{{ incrementText }}</div></div>',
  data: {
    counterVal: 0,
    incrementText: ""
  },
  created: function() {
    var vm = this;
    goui.onMessage = function(messageType, message){
      switch(messageType) {
        case "getIncText":
          vm.incrementText = message
          break;
        case "add":
          vm.counterVal = message
          break;          
      }      
    }
  },
  methods: {
    increment: function() {      
      var prevVal = this.counterVal 
      goui.add(this.counterVal, 1);
      goui.getIncText(prevVal, 1);      
    },
  }
});
