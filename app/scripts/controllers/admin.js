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

  $scope.update = function(idx) {
    var account = $scope.accounts[idx];
    var title = prompt("Enter a new title", account.title);
    if(title === null) {
      return;
    }
    var author = prompt("Enter a new author", account.author);
    if(author === null) {
      return;
    }
    account.title = title;
    account.author = author;
    account.$save();

    $scope.list(idx);
  };

}]);
