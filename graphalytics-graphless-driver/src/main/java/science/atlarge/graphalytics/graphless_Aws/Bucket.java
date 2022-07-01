package science.atlarge.graphalytics.graphless_Aws;

import com.amazonaws.AmazonServiceException;
import com.amazonaws.SdkClientException;
import com.amazonaws.auth.profile.ProfileCredentialsProvider;
import com.amazonaws.regions.Regions;
import com.amazonaws.services.s3.AmazonS3;
import com.amazonaws.services.s3.AmazonS3ClientBuilder;
import com.amazonaws.services.s3.model.DeleteObjectRequest;
import com.amazonaws.services.s3.model.GetObjectRequest;
import com.amazonaws.services.s3.model.ObjectListing;
import com.amazonaws.services.s3.model.S3Object;
import com.amazonaws.services.s3.model.S3ObjectSummary;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;

import static java.nio.charset.StandardCharsets.UTF_8;

/**
 * High-level wrapper for S3 buckets.
 *
 * @author jmasic
 */
public class Bucket {

    public static final String BUCKET_NAME = "graphless-graph-file-bucket";
    public static final Regions REGION = Regions.US_EAST_2;
    private final AmazonS3 s3;


    public Bucket() {
        s3 = AmazonS3ClientBuilder.standard()
                .withCredentials(new ProfileCredentialsProvider())
                .withRegion(REGION)
                .build();
    }

    public List<String> list() {
        ObjectListing listing = s3.listObjects(BUCKET_NAME);
        List<String> names = new ArrayList<>();

        names.addAll(listing.getObjectSummaries().stream().map(S3ObjectSummary::getKey).toList());
        while (listing.isTruncated()) {
            names.addAll(listing.getObjectSummaries().stream().map(S3ObjectSummary::getKey).toList());
        }

        return names;
    }

    public List<String> getObjectLines(String key) {
        S3Object object = s3.getObject(new GetObjectRequest(BUCKET_NAME, key));
        try (InputStreamReader objectData = new InputStreamReader(object.getObjectContent(), StandardCharsets.UTF_8)) {
            final BufferedReader reader = new BufferedReader(objectData);
            return reader.lines().toList();
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    public void delete(String fileName) {
        try {
            s3.deleteObject(new DeleteObjectRequest(BUCKET_NAME, fileName));
        } catch (AmazonServiceException e) {
            throw new RuntimeException("Error while deleting file", e);
        } catch (SdkClientException e) {
            throw new IllegalStateException("Something wrong happened while communicating with Amazon S3", e);
        }
    }
}
