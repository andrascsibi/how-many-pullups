angular.module('pullApp')

.controller('AdminCtrl', ['$scope', 'Account', function($scope, Account) {
  $scope.selected = null;

  $scope.list = function(idx){
    Account.query(function(data){
      $scope.accounts = data;
      if(idx !== undefined) {
        $scope.selected = $scope.accounts[idx];
        $scope.selected.idx = idx;
      }
    }, function(error){
      alert(error.data.error); // TODO
    });
  };

  $scope.list();

  $scope.get = function(idx){
    Account.get({id: $scope.accounts[idx].ID}, function(data){
      $scope.selected = data;
      $scope.selected.idx = idx;
    });
  };

  $scope.update = function(idx) {
  };

}]);
