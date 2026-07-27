// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3

const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    const parsed = tosca.parseComparisonArguments(currentPropertyValue, arguments);
    if (!parsed) {
        return false;
    }
    
    return tosca.compare(parsed.val1, parsed.val2) < 0;
};