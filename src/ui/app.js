var vm = new Vue({
  el: '#app',
  template: '<div><div class="counter">{{ counterVal }}</div><button class="btn btn-primary" v-on:click="increment">Increment</button></div>',
  data: {
    counterVal: 0
  },
  methods: {
    increment: function() { 
      counter.add(1); 
      this.counterVal = counter.Value;
    },
  }
});
