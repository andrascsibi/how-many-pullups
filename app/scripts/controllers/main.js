angular.module('pullApp')

.value('reloadInterval', 5 * 60 * 1000)

// .controller('TotalCtrl', ['$scope', '$http', '$interval', 'reloadInterval',
//   function($scope, $http, $interval, reloadInterval) {

//   if (!$scope.refresh) {
//     var refresh = function() {
//       $http({method: 'GET', url: 'total'}).
//         success(function(data, status, headers, config) {
//           $scope.stat = data;
//         }).
//         error(function(data, status, headers, config) {
//           console.log("request failed");
//       });
//     };
//     $scope.refresh = refresh;
//     refresh();
//     $interval(refresh, reloadInterval);
//     $scope.repButtons = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20];
//   }
// }])

// .controller('HelloCtrl', ['$scope', '$http', function($scope, $http) {
//   $http({method: 'GET', url: 'whoami'}).
//     success(function(data, status, headers, config) {
//       $scope.stat = data;
//     }).
//     error(function(data, status, headers, config) {
//       console.log("request failed");
//   });
// }])


.controller('NavbarCtrl', ['$scope', '$http', '$modal', '$resource', function($scope, $http, $modal, $resource){
  $http.get('whoami').
  success(function(data, status, headers, config) {
    $scope.whoami = data;
    if (data.Unregistered) {
      $scope.showRegModal();
    }
  }).
  error(function(data, status, headers, config) {
    console.log("request failed");
  });

  var regModal = $modal({scope: $scope, template: 'app/views/registration.html', show: false});
  $scope.showRegModal = function() {
    regModal.$promise.then(regModal.show);
  };

}])

.controller('AdminCtrl', ['$scope', '$resource', function($scope, $resource) {
  var Account = $resource("/accounts/:id", {id: '@id'}, {});

  $scope.selected = null;

  $scope.list = function(idx){
    Account.query(function(data){
      $scope.accounts = data;
      if(idx !== undefined) {
        $scope.selected = $scope.accounts[idx];
        $scope.selected.idx = idx;
      }
    }, function(error){
      alert(error.data);
    });
  };

  $scope.list();

  $scope.get = function(idx){
    Account.get({id: $scope.accounts[idx].ID}, function(data){
      $scope.selected = data;
      $scope.selected.idx = idx;
    });
  };
}])

.controller('BoardCtrl', ['$scope', '$resource', '$routeParams', function($scope, $resource, $routeParams) {
  var Account = $resource("/accounts/:id", {id: '@id'}, {});

  Account.get({id: $routeParams.id}, function(data){
    $scope.account = data;
  });
}]);

