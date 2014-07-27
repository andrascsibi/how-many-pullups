angular.module('pullApp')

.controller('ProfileCtrl', ['$scope', '$http',
  function($scope, $http) {

  var following = $scope.whoami.Account.Following;

  $scope.following = following !== null &&
      following.indexOf($scope.account.ID) >= 0;

  $scope.follow = function() {
    var follower = $scope.whoami.Account.ID;
    var followee = $scope.account.ID;
    var op = $scope.following ? 'unfollow' : 'follow';

    $scope.working = true;

    $http.post('/' + ['accounts', op, follower, followee].join('/'))
    .success(function(data, status) {
      $scope.following = !$scope.following;
      $scope.working = false;
    })
    .error(function(data, status) {
      $scope.working = false;
    });
  };

}]);
