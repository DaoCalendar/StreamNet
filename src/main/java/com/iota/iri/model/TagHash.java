package com.iota.iri.model;

/**
 * TagHash class.
 */
public class TagHash extends AbstractHash {

    protected TagHash(byte[] tagBytes, int offset, int tagSizeInBytes) {
        super(tagBytes, offset, tagSizeInBytes);
    }
}
