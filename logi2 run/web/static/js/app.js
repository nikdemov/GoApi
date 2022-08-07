angular.module('logi2', ["ngMaterial"])
    .config(["$qProvider", function($qProvider) {
        $qProvider.errorOnUnhandledRejections(false);
    }]);

angular.module('logi2', ['ngMaterial'])
    .config(function($mdThemingProvider) {
        $mdThemingProvider.theme('default')
            .primaryPalette('orange')
            .accentPalette('amber');
        $mdThemingProvider.setDefaultTheme('default');
    });