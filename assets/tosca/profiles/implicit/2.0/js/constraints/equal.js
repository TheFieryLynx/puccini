// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3

const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    const parsed = tosca.parseComparisonArguments(currentPropertyValue, arguments);
    if (!parsed) {
        return false;
    }
    
    const comparable1 = tosca.getComparable(parsed.val1);
    const comparable2 = tosca.getComparable(parsed.val2);
    
    return comparable1 === comparable2;
};