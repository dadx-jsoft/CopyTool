/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
package copy;

import java.io.File;
import org.apache.commons.io.FileUtils;

/**
 *
 * @author xuand
 */
public class CopyFilesUtils {

    private static long NUMBER_OF_FILES = 0;

//    public static void main(String[] args) {
//        String sourceDir = "D:\\origin-folder-1";
//        String destDir = "D:\\origin-folder-2";
//        CopyFilesUtils.copyFilesToDirectory(sourceDir, destDir);
//    }

    public static long copyFilesToDirectory(final String sourceDir, final String destDir) {
        File source = new File(sourceDir);
        try {
            File[] listFiles = source.listFiles();
            for (final File file : listFiles) {
                if (file.isDirectory()) {
                    copyFilesToDirectory(file.getAbsolutePath(), destDir);
                } else {
                    File dest = new File(destDir + "/" + file.getName());
                    if (dest.exists()) {
                        File duplicateFile = new File(destDir + "/" + System.currentTimeMillis() + "-" + file.getName());
                        FileUtils.copyFile(file, duplicateFile);
                    } else {
                        FileUtils.copyFile(file, dest);
                    }
                    NUMBER_OF_FILES++;
                }
            }
            return NUMBER_OF_FILES;
        } catch (Exception ex) {
            return -1;
        }
    }

    public static void resetNumberOfFiles() {
        NUMBER_OF_FILES = 0;
    }

}
