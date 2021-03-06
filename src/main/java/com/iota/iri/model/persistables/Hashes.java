package com.iota.iri.model.persistables;

import com.iota.iri.model.Hash;
import com.iota.iri.model.HashFactory;
import com.iota.iri.storage.Persistable;
import org.apache.commons.lang3.ArrayUtils;

import java.util.LinkedHashSet;
import java.util.Set;

/**
 * Hash class for storage.
 */
public class Hashes implements Persistable {
    public Set<Hash> set = new LinkedHashSet<>();
    private static final byte delimiter = ",".getBytes()[0];

    public byte[] bytes() {
        return set.parallelStream()
                .map(Hash::bytes)
                .reduce((a,b) -> ArrayUtils.addAll(ArrayUtils.add(a, delimiter), b))
                .orElse(new byte[0]);
    }

    public void read(byte[] bytes) {
        if(bytes != null) {
            set = new LinkedHashSet<>(bytes.length / (1 + Hash.SIZE_IN_BYTES) + 1);
            for (int i = 0; i < bytes.length; i += 1 + Hash.SIZE_IN_BYTES) {
                set.add(HashFactory.GENERIC.create(this.getClass(), bytes, i, Hash.SIZE_IN_BYTES));
            }
        }
    }

    @Override
    public byte[] metadata() {
        return new byte[0];
    }

    @Override
    public void readMetadata(byte[] bytes) {

    }

    @Override
    public boolean merge() {
        return true;
    }
}
