angular.module('pullApp')

.controller('BoardCtrl', ['$scope', 'Account', 'Challenge', '$routeParams', 'WhoamiService', function($scope, Account, Challenge, $routeParams, WhoamiService) {

  WhoamiService().then(function(whoami) {
    $scope.whoami = whoami;
    $scope.owner = $routeParams.id === whoami.Account.ID;
  });

  Account.get({id: $routeParams.id}, function(data){
    $scope.account = data;
  }, function(err) {
    if (err.status === 404) {
      $scope.notFound = true;
    }
  });

  $scope.challenges = [];

  $scope.list = function(idx) {
    Challenge.query({id: $routeParams.id},function(data){
      $scope.challenges = data;
    }, function(error){
      alert(error.data.error); // TODO
    });
  };

  $scope.list();

  $scope.add = function() {
    var newChallenge = new Challenge();
    newChallenge.Title = 'Pullups';
    newChallenge.Description = 'My brand new challenge for the month';
    newChallenge.$save({id: $routeParams.id})
    .then($scope.list);
  };

}]);

