angular.module('pullApp')

.controller('HomepageCtrl', ['$scope', 'WhoamiService',
  function($scope, WhoamiService){

  WhoamiService().then(function(whoami) {
    $scope.whoami = whoami;
  });
}]);
