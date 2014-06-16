angular.module('pullApp')

.controller('ChallengeCtrl', ['$scope', '$rootScope', '$interval', 'challenge', 'allSets', 'whoami',
  function($scope, $rootScope, $interval, challenge, allSets, whoami) {

  $scope.challenge = challenge;
  $rootScope.title = challenge.AccountID + "'s " + challenge.Title + ' challenge';
  $scope.whoami = whoami;

  if (allSets.length === 0) {
    $scope.empty = true;
    return;
  }

  $scope.refresh = function(c, newSet) {
    allSets.unshift(newSet);
    process(allSets);
    hourlyCal.update($scope.cal.data);
    dailyCal.update($scope.cal.data);
  };

  var process = function(allSets) {

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

    var dayKey = function(set) {
      var d = new Date(set.ts*1000);
      return new Date(d.getFullYear(), d.getMonth(), d.getDate(), 0, 0, 0, 0)
          .toLocaleString();
    };
    var hourKey = function(set) {
      var d = new Date(set.ts*1000);
      return new Date(d.getFullYear(), d.getMonth(), d.getDate(), d.getHours(), 0, 0, 0)
          .toLocaleString();
    };

    var histogram = function(sets, keyFun) {
      var groups = _.groupBy(sets, keyFun);
      var repCount = _.values(groups)
          .map(function(g) {
            return g.reduce(function(sum, cur) {
              return sum + cur.reps;
            }, 0);
          });
      return _.sortBy(repCount, function(l,r) {return l-r;});
    };

    var getPercentile = function(arr, p) {
      arr.sort(function(l,r) {return l-r;});
      return arr[Math.floor(arr.length * p)];
    };

    var getLegend = function(max) {
      var step = max / 4;
      return [1, 2, 3, 4].map(function(e){
        return Math.round(e * step);
      });
    };

    var sets = parseDates(allSets);
    var repsByDay = histogram(sets, dayKey);
    var repsByHour = histogram(sets, hourKey);

    $scope.stats = getStats(allSets);
    $scope.stats.workDays = repsByDay.length;

    $scope.cal = {
      data: toCalHeatmap(sets),
      hourlyLegend: getLegend(getPercentile(repsByHour, 0.9)),
      dailyLegend: getLegend(getPercentile(repsByDay, 0.9)),
    };
  };

  process(allSets);


  var minDate = new Date($scope.stats.minDate);
  var maxDate = new Date($scope.stats.maxDate);
  var now = new Date();
  var sixWeeksAgo = new Date().setDate(now.getDate() - 5*7);
  var isOld = sixWeeksAgo > maxDate;
  var isNew = minDate > sixWeeksAgo;

  var calSettings = {
    start: isOld || isNew ? minDate : sixWeeksAgo,
    minDate: minDate,
    maxDate: isOld ? maxDate: now,
    range: 6,
    domain: "week",
    itemName: ['rep', 'reps'],
    tooltip: true,
    data: $scope.cal.data,
    cellSize: 15,
    legendCellSize: 15,
    domainGutter: 10,
    onMinDomainReached: function(reached) {
      $scope.prevDisabled = reached;
    },
    onMaxDomainReached: function(reached) {
      $scope.nextDisabled = reached;
    },
  };

  var hourlyCal = new CalHeatMap();
  hourlyCal.init(angular.extend({
    itemSelector: '#cal-heatmap-hourly',
    itemNamespace: 'hourly',
    rowLimit: 24,
    subDomain: "hour",
    legend: $scope.cal.hourlyLegend,
    legendHorizontalPosition: 'right',
    legendVerticalPosition: 'top',
    label: {position: 'top'},
    highlight: now,
  }, calSettings));

  var dailyCal = new CalHeatMap();
  dailyCal.init(angular.extend({
    itemSelector: '#cal-heatmap-daily',
    itemNamespace: 'daily',
    rowLimit: 1,
    subDomain: "day",
    subDomainTextFormat: function(date, value) {
      return value;
    },
    legend: $scope.cal.dailyLegend,
    legendHorizontalPosition: 'right',
    legendVerticalPosition: 'bottom',
    domainLabelFormat: '',
  }, calSettings));


  var STEPS = 2;
  $scope.next = function() {
    hourlyCal.next(STEPS);
    dailyCal.next(STEPS);
  };

  $scope.previous = function() {
    hourlyCal.previous(STEPS);
    dailyCal.previous(STEPS);
  };

  $interval(function() {
    hourlyCal.highlight(new Date());
  }, 60 * 1000);

}]);
