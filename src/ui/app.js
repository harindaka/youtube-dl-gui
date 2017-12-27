//alert(document.body.getElementsByTagName('div')[0].parentElement.outerHTML);

// var vm = new Vue({
//   el: '#app',
//   template: '#counter-template',
//   data: {
//     counterVal: 0,
//     incrementText: ""
//   },
//   created: function() {

//   },
//   methods: {
//     increment: function() {        
//       var vm = this;    
//       var prevVal = this.counterVal; 
//       goui.send("add", {
//         val1: this.counterVal, 
//         val2: 1
//       }, function(result){
//         vm.counterVal = parseInt(result);
//       });
      
//       goui.send("getIncText", {
//         val1: prevVal, 
//         val2:1
//       }, function(result){
//         vm.incrementText = result;
//       });      
//     },
//   }
// });

new Vue({
  el: '#app',
  template: '#home-template',
  data: {
    websiteUrl: "https://youtube.com",
    websiteWindow: null
  },
  created: function() {
    document.domain = "youtube.com";
  },
  methods: {
    openWebsiteWindow: function() {         
      if(this.websiteWindow){
        this.websiteWindow.close();
      }
      
      this.websiteWindow = window.open(this.websiteUrl, '_blank', 'location=yes,height=570,width=520,scrollbars=yes,status=yes');            
    },

    downloadSelection: function(){
      if(this.websiteWindow){
        alert(this.websiteWindow.location);
      }
    }
  }
});
