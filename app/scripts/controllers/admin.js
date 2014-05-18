angular.module('pullApp')

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
