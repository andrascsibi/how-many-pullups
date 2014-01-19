angular.module('pullApp')
 
.value('reloadInterval', 5 * 60 * 1000)
 
.controller('TotalCtrl', ['$scope', '$http', '$interval', 'reloadInterval',
  function($scope, $http, $interval, reloadInterval) {

  if (!$scope.refresh) {
    var refresh = function() {
      $http({method: 'GET', url: 'total'}).
        success(function(data, status, headers, config) {
          $scope.stat = data;
        }).
        error(function(data, status, headers, config) {
          console.log("request failed");
      });
    };
    $scope.refresh = refresh;
    refresh();
    $interval(refresh, reloadInterval);
  }
}])
 
.controller('HelloCtrl', ['$scope', '$http', function($scope, $http) {
  $http({method: 'GET', url: 'whoami'}).
    success(function(data, status, headers, config) {
      $scope.stat = data;
    }).
    error(function(data, status, headers, config) {
      console.log("request failed");
  });
}]);
 
 
