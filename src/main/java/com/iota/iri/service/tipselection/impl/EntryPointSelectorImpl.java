package com.iota.iri.service.tipselection.impl;

import com.iota.iri.validator.MilestoneTracker;
import com.iota.iri.controllers.MilestoneViewModel;
import com.iota.iri.model.Hash;
import com.iota.iri.service.snapshot.SnapshotProvider;
import com.iota.iri.service.tipselection.EntryPointSelector;
import com.iota.iri.storage.Tangle;

/**
 * Implementation of <tt>EntryPointSelector</tt> that given a depth N, returns a N-deep milestone.
 * Meaning <CODE>milestone(latestSolid - depth)</CODE>
 * Used to as a starting point for the random walk.
 */
public class EntryPointSelectorImpl implements EntryPointSelector {

    private final Tangle tangle;
    private final SnapshotProvider snapshotProvider;
    private final MilestoneTracker milestoneTracker;

    public EntryPointSelectorImpl(Tangle tangle, SnapshotProvider snapshotProvider, MilestoneTracker milestoneTracker) {
        this.tangle = tangle;
        this.snapshotProvider = snapshotProvider;
        this.milestoneTracker = milestoneTracker;
    }

    @Override
    public Hash getEntryPoint(int depth) throws Exception {
        int milestoneIndex = Math.max(milestoneTracker.latestSolidSubtangleMilestoneIndex - depth - 1,
                snapshotProvider.getInitialSnapshot().getIndex());
        MilestoneViewModel milestoneViewModel = MilestoneViewModel.findClosestNextMilestone(tangle, milestoneIndex,
                milestoneTracker.latestMilestoneIndex);
        if (milestoneViewModel != null && milestoneViewModel.getHash() != null) {
            return milestoneViewModel.getHash();
        }

        return milestoneTracker.latestSolidSubtangleMilestone;
    }
}
