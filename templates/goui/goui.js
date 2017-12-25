(function(){
    if(!window.goui){
        window.goui = {}
    }

    goui.toES5MultilineString = function(f) {
        return f.toString().split('\n').slice(1, -1).join('\n');
    }
    
    goui.messageHandlers = {};
    
    goui.appendHtmlTemplate = function(id, templateHtml){        
        // var el = document.createElement('script', {
        //     type: 'text/x-template',
        //     id: id
        // });         
        // el.innerHTML = templateHtml;        
        // document.body.appendChild(el);

        var el = document.createElement('script', {
            id: id
        });         
        el.type = 'text/x-template';
        el.innerHTML = templateHtml;        
        //el.style.cssText = "display:none";        
        document.body.appendChild(el);      
    };

    goui.onMessage = function(messageType, messageHandler){
        goui.messageHandlers[messageType] = messageHandler;
    };

    goui.goToJs = function(messageType, message, callbackId){
        if(callbackId !== "" && goui.messageHandlers[messageType] && goui.messageHandlers[messageType][callbackId]){
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
        
        var callbackId = null;
        if(callback){
            if(!goui.messageHandlers[messageType]){
                goui.messageHandlers[messageType] = {};
            }

            callbackId = guid();
            while(goui.messageHandlers[messageType][callbackId]){
                callbackId = guid();
            }
                        
            goui.messageHandlers[messageType][callbackId] = callback;
        }
        
        if(!callbackId){
            callbackId = "";
        }

        goui.jsToGo(messageType, stringifiedMessage, callbackId);        
    }

    function guid() {
        function s4() {
          return Math.floor((1 + Math.random()) * 0x10000)
            .toString(16)
            .substring(1);
        }
        return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
          s4() + '-' + s4() + s4() + s4();
    }
})();