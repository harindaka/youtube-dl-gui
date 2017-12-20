(function(){
    // if(!window.goui){
    //     window.goui = {}
    // }

    goui.messageHandlers = {};
    goui.onMessage = function(messageType, messageHandler){
        goui.messageHandlers[messageType] = messageHandler;
    };

    goui.invokeJsMessageHandler = function(messageType, message, callbackId){
        var handler = null;
        if(goui.messageHandlers[messageType] && goui.messageHandlers[messageType][callbackId]){
            var handler = goui.messageHandlers[messageType][callbackId];
            delete goui.messageHandlers[messageType][callbackId];
            var parsedMessage = JSON.parse(message);            
            handler(parsedMessage);
        }        
    };

    goui.send = function(messageType, messageObject, callback){
        var stringifiedMessage = "";
        if(typeof messageObject !== 'undefined' && messageObject !== null){
            stringifiedMessage = JSON.stringify(messageObject);
        }

        if(callback){
            if(!goui.messageHandlers[messageType]){
                goui.messageHandlers[messageType] = {};
            }

            var callbackId = uuidv4();
            while(goui.messageHandlers[messageType][callbackId]){
                callbackId = uuidv4();
            }

            goui.messageHandlers[messageType][callbackId] = callback;
        }

        


    var res = [];
    for(var m in goui) {
        if(typeof goui[m] == "function") {
            res.push(m)
        }
    }
    alert(res);




        goui.invokeGoMessageHandler(messageType, stringifiedMessage, callbackId);        
    }

    function uuidv4() {
        // return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, function(c){
        //   (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
        // });
        function guid() {
            function s4() {
              return Math.floor((1 + Math.random()) * 0x10000)
                .toString(16)
                .substring(1);
            }
            return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
              s4() + '-' + s4() + s4() + s4();
        }
    }
})();