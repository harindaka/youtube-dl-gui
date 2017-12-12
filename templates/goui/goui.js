(function(){
    if(!window.goui){
        window.goui = {}
    }

    goui.messageHandlers = {};
    goui.onMessage = function(messageType, messageHandler){
        goui.messageHandlers[messageType] = messageHandler;
    };

    goui.invokeJsMessageHandler = function(messageType, message){
        var handler = goui.messageHandlers[messageType];
        if(handler){
            var parsedMessage = JSON.parse(message)
            handler(parsedMessage);
        }
    };

    goui.send = function(messageType, messageObject){
        var stringifiedMessage = "";
        if(typeof messageObject !== 'undefined' && messageObject !== null){
            stringifiedMessage = JSON.stringify(messageObject);
        }

        goui.invokeGoMessageHandler(messageType, stringifiedMessage);
    }
})();