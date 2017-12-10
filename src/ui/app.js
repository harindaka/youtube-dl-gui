var vm = new Vue({
  el: '#app',
  template: '<div><div class="counter">{{ counterVal }}</div><button class="btn btn-primary" v-on:click="increment">Increment</button><div>{{ incrementText }}</div></div>',
  data: {
    counterVal: 0,
    incrementText: ""
  },
  created: function() {
    var vm = this;
    native.done = function(method, result){
      switch(method) {
        case "getIncText":
          vm.incrementText = result
          break;
        case "add":
          vm.counterVal = result
          break;          
      } 
      
    }
  },
  methods: {
    increment: function() { 
      var prevVal = this.counterVal
      native.add(this.counterVal, 1);  
      native.getIncText(prevVal, 1);
    },
  }
});
