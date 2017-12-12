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

    goui.send = function(messageType, message){
        var stringifiedMessage = "";
        if(typeof message !== 'undefined' && message !== null){
            stringifiedMessage = JSON.stringify(message);
        }

        goui.invokeGoMessageHandler(messageType, stringifiedMessage);
    }
})();