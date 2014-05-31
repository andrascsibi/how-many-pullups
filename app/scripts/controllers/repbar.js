angular.module('pullApp')

.controller('RepbarCtrl', ['$scope', 'Sets', function($scope, Sets) {
  $scope.repButtons = [1,2,3,4,5,6,7,8,9,10];

  $scope.add = function(reps) {

    $scope.working = true;
    var newSet = new Sets();
    newSet.Reps = reps;
    newSet.$save({id: $scope.challenge.AccountID, c_id: $scope.challenge.ID})
    .then(function(){
      $scope.working = false;
      $scope.list();
    }, function(err) {
      alert(err.data.error); // TODO
    });
  };

}]);

