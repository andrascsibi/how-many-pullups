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

  // $scope.reloadChallenge = function(c) {
  //   $scope.challenges.forEach(function(e, i) {
  //     if (e.ID !== c.ID) return;
  //     e.$get(function(data) {
  //       $scope.challenges = $scope.challenges.map(function(e) {
  //         if (e.ID !== c.ID) return e;
  //         return data;
  //       });
  //     });
  //   });
  // };

  $scope.list = function() {
    Challenge.query({id: $routeParams.id}, function(data){
      $scope.challenges = data;
    }, function(error){
      alert(error.data.error); // TODO
    });
  };

  $scope.list();

  $scope.edited = null;

  $scope.add = function() {
    var newChallenge = new Challenge();
    newChallenge.Title = 'Pullups';
    newChallenge.Description = '';
    newChallenge.AccountID = $routeParams.id;
    $scope.challenges.splice(0, 0, newChallenge);
    $scope.edited = newChallenge;
  };

  $scope.edit = function(c) {
    $scope.edited = angular.copy(c);
  };

  $scope.cancel = function() {
    $scope.edited = null;
    $scope.list();
  };

  $scope.save = function() {
    $scope.working = true;
    $scope.edited.$save()
    .then(function(){
    }, function(err) {
      alert(err.data.error);
    })
    .then(function(){
      $scope.working = false;
      $scope.edited = null;
      $scope.list();
    });
  };


}]);

