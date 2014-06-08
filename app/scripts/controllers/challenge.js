angular.module('pullApp')

.controller('ChallengeCtrl', ['$scope', '$rootScope', 'challenge', 'allSets', 'WhoamiService',
  function($scope, $rootScope, challenge, allSets, WhoamiService) {

  $scope.challenge = challenge;
  $rootScope.title = challenge.AccountID + "'s " + challenge.Title + ' challenge';

  var parseDates = function(sets) {
    return sets
    .filter(function(cur) {
      return cur !== 0;
    })
    .map(function(cur) {
      return {
        ts: Math.round(Date.parse(cur.Date)/1000),
        reps: cur.Reps
      };
    });
  };

  var toCalHeatmap = function(sets) {
    return sets.reduce(function(memo, cur) {
      if (!memo[cur.ts]) {
        memo[cur.ts] = 0;
      }
      memo[cur.ts] += cur.reps;
      return memo;
    }, {});
  };

  var getStats = function(sets) {
    if (sets.length === 0) {
      return null;
    }
    return sets.reduce(function(memo, cur) {
      memo.numSets++;
      memo.totalReps += cur.Reps;
      memo.maxReps = Math.max(memo.maxReps, cur.Reps);
      return memo;
    }, {
      numSets: 0,
      totalReps: 0,
      maxReps: 0,
      minDate: sets[sets.length - 1].Date,
      maxDate: sets[0].Date,
    });
  };

  var sets = parseDates(allSets);
  $scope.stats = getStats(allSets);

//  var setsByDay = groupByDates(sets);

  var cal = new CalHeatMap();
  cal.init({
    start: new Date(2014, 0), // January, 1st 2000
    range: 6,
    end: new Date(),
    itemName: ['rep', 'reps'],
    tooltip: true,
    domain: "month",
    subDomain: "day",
    data: toCalHeatmap(sets),
    legend: [10, 30, 50, 75]
  });


}]);
