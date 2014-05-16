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
    $scope.repButtons = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20];
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
}])


.controller('HomePageCtrl', ['$scope', '$http', '$location', function($scope, $http, $location){
  $http.get('whoami').
  success(function(data, status, headers, config) {
    console.log(status);
    $scope.whoami = data;
  }).
  error(function(data, status, headers, config) {
    console.log("request failed");
  });

}]);

