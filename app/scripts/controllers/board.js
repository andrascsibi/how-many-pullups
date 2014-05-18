angular.module('pullApp')

.controller('BoardCtrl', ['$scope', '$resource', '$routeParams', 'WhoamiService', function($scope, $resource, $routeParams, WhoamiService) {
  var Account = $resource("/accounts/:id", {id: '@id'}, {});
  WhoamiService().then(function(whoami) {
    $scope.whoami = whoami;
  });

  Account.get({id: $routeParams.id}, function(data){
    $scope.account = data;
  }, function(err) {
    if (err.status === 404) {
      $scope.notFound = true;
    }
  });
}]);

