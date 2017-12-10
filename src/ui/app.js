var vm = new Vue({
  el: '#app',
  template: '<div><div class="counter">{{ counterVal }}</div><button class="btn btn-primary" v-on:click="increment">Increment</button></div>',
  data: {
    counterVal: 0
  },
  methods: {
    increment: function() {    
      native.done = function(result){
        alert(result);
        this.counterVal = result
      } 
      native.add(this.counterVal, 1);                   
    },
  }
});
